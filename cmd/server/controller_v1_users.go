package main

func (controller *Controller) deleteUserByID(ID int) (any, error) {
	return nil, controller.app.userRepo.DeleteUserByID(ID)
}

func (controller *Controller) getUsers() (any, error) {
	return controller.app.userRepo.GetUsers()
}
