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

func InSession(session func(operator *Operator) error) error {
	s, err := sessions.Start()
	if err != nil {
		return fmt.Errorf("произошла ошибка при запуске сессии: %s", err)
	}

	defer s.Close()

	o := &Operator{
		s: s,
	}

	return session(o)
}
