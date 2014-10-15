func boardcast(msg interface{}) {
    // get user list
    for user := range user_list {
        user.c <- &SSysmsg{"map board", msg}
    }
}

func Dispatch(user *User, msg_type uint32, msg interface{}) {
    handle := msg_handle[msg_type]
    handle(user, msg)
}
