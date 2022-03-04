package errutil

type Convert func(error) (error, bool)

type ConvertChain []Convert

func (cc ConvertChain) Exec(err error) error {
	if err == nil {
		return nil
	}
	for _, c := range cc {
		if cErr, ok := c(err); ok {
			return cErr
		}
	}
	return err
}
