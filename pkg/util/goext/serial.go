package goext

import (
	"fmt"
	"github.com/26597925/EastCloud/pkg/util/stringext"
	"github.com/pkg/errors"
	"go.uber.org/multierr"
)

func SerialWithError(fns ...func() error) func() error {
	return func() error {
		var errs error
		for _, fn := range fns {
			errs = multierr.Append(errs, try(fn, nil))
		}
		return errs
	}
}

func SerialUntilError(fns ...func() error) func() error {
	return func() error {
		for _, fn := range fns {
			if err := try(fn, nil); err != nil {
				return errors.Wrap(err, stringext.FunctionName(fn))
			}
		}
		return nil
	}
}

func try(fn func() error, cleaner func()) (ret error) {
	if cleaner != nil {
		defer cleaner()
	}
	defer func() {
		if err := recover(); err != nil {
			//_, file, line, _ := runtime.Caller(2)
			//_logger.Error("recover", zap.Any("err", err), zap.String("line", fmt.Sprintf("%s:%d", file, line)))
			if _, ok := err.(error); ok {
				ret = err.(error)
			} else {
				ret = fmt.Errorf("%+v", err)
			}
		}
	}()
	return fn()
}

func try2(fn func(), cleaner func()) (ret error) {
	if cleaner != nil {
		defer cleaner()
	}
	defer func() {
		//_, file, line, _ := runtime.Caller(4)
		if err := recover(); err != nil {
			//_logger.Error("recover", zap.Any("err", err), zap.String("line", fmt.Sprintf("%s:%d", file, line)))
			if _, ok := err.(error); ok {
				ret = err.(error)
			} else {
				ret = fmt.Errorf("%+v", err)
			}
		}
	}()
	fn()
	return nil
}