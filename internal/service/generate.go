package service

//go:generate sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate minimock -i EventsService -o ./mocks/ -s "_minimock.go"
//go:generate minimock -i TicketsService -o ./mocks/ -s "_minimock.go"
//go:generate minimock -i UsersService -o ./mocks/ -s "_minimock.go"
