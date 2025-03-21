// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package azureeventhubreceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/azureeventhubreceiver"

import (
	"context"
	"errors"

	eventhub "github.com/Azure/azure-event-hubs-go/v3"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/extension/xextension/storage"
	"go.opentelemetry.io/collector/receiver"
	"go.uber.org/zap"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/adapter"
)

type hubWrapper interface {
	GetRuntimeInformation(ctx context.Context) (*eventhub.HubRuntimeInformation, error)
	Receive(ctx context.Context, partitionID string, handler eventhub.Handler, opts ...eventhub.ReceiveOption) (listerHandleWrapper, error)
	Close(ctx context.Context) error
}

type hubWrapperImpl struct {
	hub *eventhub.Hub
}

func (h *hubWrapperImpl) GetRuntimeInformation(ctx context.Context) (*eventhub.HubRuntimeInformation, error) {
	return h.hub.GetRuntimeInformation(ctx)
}

func (h *hubWrapperImpl) Receive(ctx context.Context, partitionID string, handler eventhub.Handler, opts ...eventhub.ReceiveOption) (listerHandleWrapper, error) {
	l, err := h.hub.Receive(ctx, partitionID, handler, opts...)
	return l, err
}

func (h *hubWrapperImpl) Close(ctx context.Context) error {
	return h.hub.Close(ctx)
}

type listerHandleWrapper interface {
	Done() <-chan struct{}
	Err() error
}

type eventhubHandler struct {
	hub           hubWrapper
	dataConsumer  dataConsumer
	config        *Config
	settings      receiver.Settings
	cancel        context.CancelFunc
	storageClient storage.Client
}

func (h *eventhubHandler) run(ctx context.Context, host component.Host) error {
	ctx, h.cancel = context.WithCancel(ctx)

	if h.storageClient == nil { // set manually for testing.
		storageClient, err := adapter.GetStorageClient(ctx, host, h.config.StorageID, h.settings.ID)
		if err != nil {
			h.settings.Logger.Debug("Error connecting to Storage", zap.Error(err))
			return err
		}
		h.storageClient = storageClient
	}

	if h.hub == nil { // set manually for testing.
		hub, newHubErr := eventhub.NewHubFromConnectionString(h.config.Connection, eventhub.HubWithOffsetPersistence(&storageCheckpointPersister{storageClient: h.storageClient}))
		if newHubErr != nil {
			h.settings.Logger.Debug("Error connecting to Event Hub", zap.Error(newHubErr))
			return newHubErr
		}
		h.hub = &hubWrapperImpl{hub: hub}
	}

	if h.config.Partition != "" {
		err := h.setUpOnePartition(ctx, h.config.Partition, true)
		if err != nil {
			h.settings.Logger.Debug("Error setting up partition", zap.Error(err))
		}
		return err
	}

	// listen to each partition of the Event Hub
	runtimeInfo, err := h.hub.GetRuntimeInformation(ctx)
	if err != nil {
		h.settings.Logger.Debug("Error getting Runtime Information", zap.Error(err))
		return err
	}

	var errs []error
	for _, partitionID := range runtimeInfo.PartitionIDs {
		err = h.setUpOnePartition(ctx, partitionID, false)
		if err != nil {
			h.settings.Logger.Debug("Error setting up partition", zap.Error(err), zap.String("partition", partitionID))
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (h *eventhubHandler) setUpOnePartition(ctx context.Context, partitionID string, applyOffset bool) error {
	receiverOptions := []eventhub.ReceiveOption{}
	if applyOffset && h.config.Offset != "" {
		receiverOptions = append(receiverOptions, eventhub.ReceiveWithStartingOffset(h.config.Offset))
	}

	if h.config.ConsumerGroup != "" {
		receiverOptions = append(receiverOptions, eventhub.ReceiveWithConsumerGroup(h.config.ConsumerGroup))
	}

	handle, err := h.hub.Receive(ctx, partitionID, h.newMessageHandler, receiverOptions...)
	if err != nil {
		return err
	}
	go func() {
		<-handle.Done()
		err := handle.Err()
		if err != nil {
			h.settings.Logger.Error("Error reported by event hub", zap.Error(err))
		}
	}()

	return nil
}

func (h *eventhubHandler) newMessageHandler(ctx context.Context, event *eventhub.Event) error {
	err := h.dataConsumer.consume(ctx, event)
	if err != nil {
		h.settings.Logger.Error("error decoding message", zap.Error(err))
		return err
	}

	return nil
}

func (h *eventhubHandler) close(ctx context.Context) error {
	var errs error
	if h.storageClient != nil {
		if err := h.storageClient.Close(ctx); err != nil {
			errs = errors.Join(errs, err)
		}
		h.storageClient = nil
	}

	if h.hub != nil {
		err := h.hub.Close(ctx)
		if err != nil {
			errs = errors.Join(errs, err)
		}
		h.hub = nil
	}
	if h.cancel != nil {
		h.cancel()
	}

	return errs
}

func (h *eventhubHandler) setDataConsumer(dataConsumer dataConsumer) {
	h.dataConsumer = dataConsumer
}

func newEventhubHandler(config *Config, settings receiver.Settings) *eventhubHandler {
	return &eventhubHandler{
		config:   config,
		settings: settings,
	}
}
