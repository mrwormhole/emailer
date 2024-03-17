# emailer

tiny email server for managed email providers such as Brevo and others to follow

### Purpose

packaging SMTP(send only) APIs to a choosable option is the goal here, then to have a tiny image that picks 1 provider and serves
as HTTP server so other microservices that belong to you can make easy email sends for multiple clients of yours.

### Supported providers

[ ] Brevo (highly recommended)
[ ] Postmark
[ ] Mailchimp
[ ] Mailtrap
[ ] Mailjet
[ ] Mailgun
[ ] Sendgrid

Note: Anything that only uses oauth2 like zoho does will not be implemented here for foreseeable future 
