package server

/* TODO (drausin) enable once have in-memory storer
func TestStart(t *testing.T) {
	up := make(chan *Key, 1)
	wg1 := new(sync.WaitGroup)
	wg1.Add(1)
	go func(wg2 *sync.WaitGroup) {
		defer wg2.Done()
		err := Start(NewDefaultConfig(), up)
		assert.Nil(t, err)
	}(wg1)

	x := <-up
	assert.NotNil(t, x)

	x.StopServer()
	wg1.Wait()
}
*/
