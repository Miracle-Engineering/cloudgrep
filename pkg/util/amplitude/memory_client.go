package amplitude

import "go.uber.org/zap"

//MemoryClient is used for testing Amplitude events are sent without sending anything to Amplitude API
type MemoryClient struct {
	UploadRequests []map[string]interface{}
	err            error
}

func UseMemoryClient() *MemoryClient {
	result := &MemoryClient{}
	defaultClient = result
	return result
}

func (c *MemoryClient) SendEvent(logger *zap.Logger, eventType EventType, eventProperties map[string]string) {
	var request map[string]interface{}
	request, c.err = newUploadRequest(eventType, eventProperties)
	if c.err == nil {
		c.UploadRequests = append(c.UploadRequests, request)
	}
}

func (c *MemoryClient) Size() int {
	return len(c.UploadRequests)
}

func (c *MemoryClient) LastEvent() (map[string]interface{}, error) {
	if c.err != nil {
		return nil, c.err
	}
	if c.Size() != 0 {
		uploadReq := c.UploadRequests[len(c.UploadRequests)-1]
		event := uploadReq["events"].([]map[string]interface{})[0]
		return event, nil
	}
	return nil, nil
}
