package group

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func Test0001(t *testing.T) {
	g := NewGroup("", nil, nil)
	g.Go("", func() error {
		fmt.Println("a")
		return nil
	})
	g.Wait()
}

func Test0002(t *testing.T) {
	g := NewGroup("", nil, nil)
	g.Go("", func() error {
		for !g.IsExit() {
			select {
			case <-g.Done():
				return nil
			case <-time.After(time.Second):
				fmt.Println("tick")
			}
		}
		return nil
	})
	g.Go("", func() error {
		time.Sleep(time.Second * 5)
		return errors.New("done")
	})
	fmt.Println(g.Wait())
}

func Test0003(t *testing.T) {
	g := NewGroup("", nil, nil)
	gg := NewGroup("", g, nil)
	gg.Go("", func() error {
		for {
			select {
			case <-gg.Done():
				return nil
			case <-time.After(time.Second):
				fmt.Println("tick")
			}
		}
		return nil
	})
	g.Go("", func() error {
		time.Sleep(time.Second * 5)
		return errors.New("done")
	})
	fmt.Println(g.Wait())
}

func Test0004(t *testing.T) {
	g := NewGroup("", nil, nil)
	gg := NewGroup("", g, nil)
	g.Go("", func() error {
		for {
			select {
			case <-g.Done():
				return nil
			case <-time.After(time.Second):
				fmt.Println("tick father")
			}
		}
		return nil
	})
	g.Go("", func() error {
		time.Sleep(time.Second * 10)
		return errors.New("done father")
	})
	gg.Go("", func() error {
		for {
			select {
			case <-gg.Done():
				return nil
			case <-time.After(time.Second):
				fmt.Println("tick")
			}
		}
		return nil
	})
	gg.Go("", func() error {
		time.Sleep(time.Second * 5)
		return errors.New("done")
	})
	fmt.Println(gg.Wait())
	fmt.Println(g.Wait())
}

func Test0005(t *testing.T) {
	g := NewGroup("", nil, func() {
		fmt.Println("stop")
	})

	g.Go("", func() error {
		for {
			select {
			case <-g.Done():
				return nil
			case <-time.After(time.Second):
				fmt.Println("tick 1")
			}
		}
		return nil
	})

	g.Go("", func() error {
		for {
			select {
			case <-g.Done():
				return nil
			case <-time.After(time.Second):
				fmt.Println("tick 2")
			}
		}
		return nil
	})

	time.Sleep(time.Second * 5)
	g.Stop()
	g.Wait()

}

func Test0006(t *testing.T) {
	done := 0
	g := NewGroup("", nil, func() {
		done = 1
		fmt.Println("stop")
	})

	g.Go("", func() error {
		for done == 0 {
			fmt.Println("tick 1")
			time.Sleep(time.Second)
		}
		return nil
	})

	g.Go("", func() error {
		for done == 0 {
			fmt.Println("tick 2")
			time.Sleep(time.Second)
		}
		return nil
	})

	time.Sleep(time.Second * 5)
	g.Stop()
	g.Wait()

}

func Test0007(t *testing.T) {
	g := NewGroup("", nil, func() {
		fmt.Println("stop")
	})

	g.Go("", func() error {
		time.Sleep(time.Second)
		fmt.Println("tick 1")
		return nil
	})

	g.Go("", func() error {
		time.Sleep(time.Second)
		time.Sleep(time.Second)
		fmt.Println("tick 2")
		return nil
	})

	g.Go("", func() error {
		for {
			select {
			case <-g.Done():
				return nil
			case <-time.After(time.Second):
				fmt.Println("tick 3")
			}
		}
		return nil
	})

	time.Sleep(time.Second * 5)
	g.Stop()
	g.Wait()
}

func Test008(t *testing.T) {
	g := NewGroup("", nil, func() {
		fmt.Println("stop")
	})

	exit := false
	g.Go("test", func() error {
		for !exit {
			fmt.Println("tick")
			time.Sleep(time.Second)
		}
		return nil
	})

	go func() {
		time.Sleep(time.Second * 5)
		g.Stop()
	}()

	go func() {
		time.Sleep(time.Second * 7)
		exit = true
	}()

	g.Wait()

}
