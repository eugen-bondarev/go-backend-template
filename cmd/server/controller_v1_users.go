package main

func (controller *Controller) deleteUserByID(ID int) error {
	return controller.app.userRepo.DeleteUserByID(ID)
}
