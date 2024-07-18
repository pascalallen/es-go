package projection

import "testing"

func TestThatNameReturnsExpectedValue(t *testing.T) {
	projection := UserEmailAddresses{}

	if projection.Name() != "user-email-addresses" {
		t.Fatal("test assertion failed for UserEmailAddresses.Name()")
	}
}

func TestThatScriptReturnsExpectedValue(t *testing.T) {
	projection := UserEmailAddresses{}
	script := `fromCategory('user')
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

	if projection.Script() != script {
		t.Fatal("test assertion failed for UpdateUserEmailAddress.CommandName()")
	}
}
