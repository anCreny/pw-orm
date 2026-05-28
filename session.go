package pworm

import (
	"fmt"

	"github.com/anCreny/pw-orm/sessions"
)

//
//  record, err := operator.NewCommandBuilder(DnsRecordGet).Where(name == test AND data == 1.1.1.1).Build().Run()
//  ...
//
//  cloneRecord := record.Clone()
//  err := cloneRecord.TrySet("RecordData.IPv4Address", "2.2.2.2")
//  ...
//
//  operator.NewCommandBuilder(DnsRecordSet).SetArguments(-NewInputObj = cloneRecord.PW(), -OldInputObj = record.PW()).Build().Run()

func NewSession() (*Operator, error) {
	s, err := sessions.Start()
	if err != nil {
		return nil, fmt.Errorf("произошла ошибка при запуске сессии: %s", err)
	}

	o := &Operator{
		s: s,
	}

	return o, nil
}

func (o *Operator) CloseSession() {
	o.s.Close()
}

func (o *Operator) SessionID() string {
	return o.s.ID
}
