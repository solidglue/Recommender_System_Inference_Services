package faiss

type FaissFactory struct {
}

func (f *FaissFactory) FaissClientFactory(indexConfStr string) *FaissIndexClient {
	ff := new(FaissIndexClient)
	ff.ConfigLoad("", "", indexConfStr)

	return ff
}
