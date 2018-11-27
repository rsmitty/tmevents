package tmevents

import (
	"net/http"
	"time"

	"github.com/cloudevents/sdk-go/v01"
)

// EventClient holds info about publishing event
type EventClient struct {
	channel   string
	namespace string
}

// NewEventClient sets info for pushing events later
func NewEventClient(channel string, namespace string) *EventClient {
	ec := EventClient{
		channel:   channel,
		namespace: namespace,
	}
	return &ec
}

// PushEvent pushes an event to a kubernetes service
func PushEvent(msgData []byte, msgID string, msgTime time.Time, ec *EventClient) error {

	//Setup event info
	event := &v01.Event{
		ContentType: "application/json",
		Data:        msgData,
		EventID:     msgID,
		EventTime:   &msgTime,
		EventType:   "cloudevent.greet.you",
		Source:      "from-galaxy-far-far-away",
	}

	//Marshal up event JSON and prepare request
	marshaller := v01.NewDefaultHTTPMarshaller()
	req, _ := http.NewRequest("POST", "http://"+ec.channel+"-channel."+ec.namespace+".svc.cluster.local", nil)
	err := marshaller.ToRequest(req, event)
	if err != nil {
		return err
	}

	//Issue POST request, but return before acking the message if there's an error
	_, err = (*http.Client).Do(&http.Client{}, req)
	if err != nil {
		return err
	}

	return nil
}
