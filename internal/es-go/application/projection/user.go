package projection

type UserEmailAddresses struct{}

func (p UserEmailAddresses) Name() string {
	return "user-email-addresses"
}

func (p UserEmailAddresses) Script() string {
	return `fromCategory('user')
.when({
    $init: function () {
        return { email_addresses: {} };
    },
    UserRegistered: function (state, event) {
        state.email_addresses[event.data.email_address] = event.data.id;
    },
    UserEmailAddressUpdated: function (state, event) {
        // Find and remove the old email associated with the id
        for (let emailAddress in state.email_addresses) {
            if (state.email_addresses[emailAddress] === event.data.id) {
                delete state.email_addresses[emailAddress];
                break;
            }
        }
        // Update to the new email
        state.email_addresses[event.data.email_address] = event.data.id;
    }
})
.outputState();`
}
