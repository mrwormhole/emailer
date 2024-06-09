# emailer

tiny email server for managed email providers

### Purpose

packaging SMTP(send only) APIs to a choosable option is the goal here, then to have a tiny image that picks 1 provider and serves
as HTTP server so other microservices that belong to you can make email sending for multiple clients of yours.

### Supported providers

- [X] Brevo (highly recommended)
- [ ] Resend
- [ ] Postmark
- [ ] Mailchimp
- [ ] Mailtrap
- [ ] Mailjet
- [ ] Mailgun
- [ ] Sendgrid (avoid them if you can)

Note: Anything that only uses oauth2 like zoho does will not be implemented here for foreseeable future 
