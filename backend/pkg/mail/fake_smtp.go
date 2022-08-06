package mail

type FakeSMTP struct{}

func NewFakeSMTP() *FakeSMTP {
	return &FakeSMTP{}
}

func (receiver *FakeSMTP) Send(input SendEmailInput) error {
	return nil
}
