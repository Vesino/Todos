package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *application) serve() error {
	// Declare a HTTP server using the same settings as in our main() function
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	shutDownError := make(chan error)
	// start a background go routine
	go listenSignals(app, srv, shutDownError)

	app.logger.PrintInfo("starting server", map[string]string{
		"addr": srv.Addr,
		"env":  app.config.env,
	})

	// start the server as normal
	err := srv.ListenAndServe()

	// Calling Shutdown() on our server will cause ListenAndServe() to immediately
	// return a http.ErrServerClosed error. So if we see this error, it is actually app
	// good thing and an indication that the graceful shutdown has started. So we check
	// specifically for this, only returning the error if it is NOT http.ErrServerClosed.
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	// Otherwise, we wait to receive the return value from Shutdown() on the
	// shutdownError channel. If return value is an error, we know that there was a
	// problem with the graceful shutdown and we return the error
	err = <-shutDownError
	if err != nil {
		return err
	}

	// at this point we know that the graceful shutdown completed successfully and we
	// log a "stopped server" message
	app.logger.PrintInfo("stopped server", map[string]string{
		"addr": srv.Addr,
	})
	return nil

}

func listenSignals(app *application, srv *http.Server, errorChan chan error) {
	// create a quit channel wich carries os.Signal values
	quit := make(chan os.Signal, 1)

	// Use signal.Notify to listen for incoming SIGINT adn SIGTERM signals and
	// relay them to the quit channel
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Read the signal from the quit channel, this block will block until a signal is received
	// This happens because the nature of the channels communication
	// it blocks until the channel can release a value, on a buffered channel, the channel can act as
	// a queue, but on a unbuffered channel we can sync the channel's communication
	s := <-quit

	// log a message to say tha the signal has been caught.
	app.logger.PrintInfo("shutting down server", map[string]string{
		"signal": s.String(),
	})

	// create a context with 5 secont time out
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Call Shutdown() on our server, passing in the context we just made.
	// Shutdown() will return nil if the graceful shutdown was successful, or an
	// error (which may happen because of a problem closing the listeners, or
	// because the shutdown didn't complete before the 5-second context deadline is
	// hit). We relay this return value to the shutdownError channel.
	err := srv.Shutdown(ctx)
	if err != nil {
		errorChan <- err
	}

	app.logger.PrintInfo("completing backround tasks", map[string]string{
		"addr": srv.Addr,
	})

	app.wg.Wait()
	errorChan <- nil
}
