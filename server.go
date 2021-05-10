package workgroup

func Server(serve func() error, shutdown func() error) RunFunc {
	return func(stop <-chan struct{}) error {
		done := make(chan error)
		defer close(done)

		go func() {
			done <- serve()
		}()

		select {
		case err := <-done:
			return err
		case <-stop:
			err := shutdown()
			if err == nil {
				err = <-done
			} else {
				<-done
			}
			return err
		}
	}
}
