A GO wrapper for Tuya's mobile app [API](https://docs.tuya.com/en/cloudapi/appAPI/index.html), based on [TuyaAIP cloud](https://github.com/TuyaAPI/cloud) and [Tuya sign hacking](https://github.com/nalajcie/tuya-sign-hacking).

```go
	t, _ := tuya.New(tuya.Config{
		CertSign:  "",
		Secret:    "",
		AppSecret: "",
	}, "custom device id", "")
	res := t.Call("tuya.m.country.list", "2.0", nil, false)
```