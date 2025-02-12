# Emailer

[![Version](https://img.shields.io/github/tag/mrwormhole/emailer.svg)](https://github.com/mrwormhole/emailer/tags)
[![CI Build](https://github.com/mrwormhole/emailer/actions/workflows/test.yaml/badge.svg)](https://github.com/mrwormhole/emailer/actions/workflows/test.yaml)
[![GoDoc](https://godoc.org/github.com/mrwormhole/emailer?status.svg)](https://godoc.org/github.com/mrwormhole/emailer)
[![Report Card](https://goreportcard.com/badge/github.com/mrwormhole/emailer)](https://goreportcard.com/report/github.com/mrwormhole/emailer)
[![License](https://img.shields.io/github/license/mrwormhole/emailer)](https://github.com/mrwormhole/emailer/blob/main/LICENSE)
[![Coverage Status](https://coveralls.io/repos/github/mrwormhole/emailer/badge.svg?branch=main)](https://coveralls.io/github/mrwormhole/emailer?branch=main)

### Purpose

Packaging SMTP(send only) APIs to a choosable option is the goal here, then to have a tiny server that picks 1 provider and serves
as HTTP endpoint so the other microservices that belong to you can make email requests.

### Supported providers

- [X] Brevo (highly recommended)
- [X] Resend
- [ ] Postmark
- [ ] Mailchimp
- [ ] Mailtrap
- [ ] Mailjet
- [ ] Mailgun
- [ ] Sendgrid (avoid them if you can)

Note: Anything that only uses oauth2 like zoho does will not be implemented here for foreseeable future 
