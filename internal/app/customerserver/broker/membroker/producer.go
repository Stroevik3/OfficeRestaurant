package membroker

import "github.com/Shopify/sarama"

type memProducer struct {
	sarama.SyncProducer
	Msg *sarama.ProducerMessage
}

func (p *memProducer) SendMessage(msg *sarama.ProducerMessage) (partition int32, offset int64, err error) {
	p.Msg = msg
	return partition, offset, err
}

func Create() *memProducer {
	p := &memProducer{}
	return p
}

func IniProducer(mp *memProducer) sarama.SyncProducer {
	return mp
}
