package controller

//PostHandler -
func PostHandler(body []byte){
	SendToQueue(body)
}

