package persistopt

import (
	"Orca_Puppet/define/debug"
	"fmt"
	"golang.org/x/sys/windows/svc/eventlog"
	"golang.org/x/sys/windows/svc/mgr"
	"strings"
)

func AddWinSvc(name, path, args, desc string, started bool) error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	s, err := m.OpenService(name)
	if err == nil {
		s.Close()
		return fmt.Errorf("service %s already exists", name)
	}
	config := mgr.Config{
		StartType:   mgr.StartAutomatic,
		DisplayName: name,
		Description: desc,
	}
	param := strings.Split(args, " ")
	s, err = m.CreateService(name, path, config, param...)
	if err != nil {
		return err
	}
	defer s.Close()
	err = eventlog.InstallAsEventCreate(name, eventlog.Error|eventlog.Warning|eventlog.Info)
	if err != nil {
		s.Delete()
		return fmt.Errorf("SetupEventLogSource() failed: %s", err)
	}
	debug.DebugPrint("service install successfully")
	if started {
		go func() {
			err = startService(name)
		}()
		if err != nil {
			return err
		}
	}
	return nil
}

func startService(name string) error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	s, err := m.OpenService(name)
	if err != nil {
		return fmt.Errorf("could not access service: %v", err)
	}
	defer s.Close()
	err = s.Start("")
	if err != nil {
		return fmt.Errorf("could not start service: %v", err)
	}
	return nil
}
