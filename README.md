# golatch1
Eleven paths Latch minimum client in golang

#Example

Import the library

    import "github.com/amassanet/golatch1"

Initialize with the application id and the secret

    la := golatch1.NewLatchApp(appId, appSecret) 

Or initialize  with custom transport 

    la := golatch1.NewLatchAppWithTransport(appId, appSecret,&http.Client{})

Pair and create an account

    accountId, err := la.Pair(token)

Get the status

    statusIsOn, err = la.StatusIsOn(accountId);

Unpair the account

    err := la.Unpair(accountId);
