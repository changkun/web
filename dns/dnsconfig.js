// Copyright 2026 Changkun Ou. All rights reserved.
// Use of this source code is governed by a GPL-3.0
// license that can be found in the LICENSE file.

// The DNS of the domains on Cloudflare, applied with DNSControl. A push
// makes each zone match this file: records missing here are deleted, and
// so are single redirect rules. Run `dnscontrol preview` first.

var REG_NONE = NewRegistrar("none");
var DSP_CLOUDFLARE = NewDnsProvider("cloudflare", {
  manage_single_redirects: true,
});

// The droplet behind traefik.
var SERVER = "138.68.107.228";

// Mail is Google Workspace on every domain.
var GOOGLE_MAIL = [
  MX("@", 1, "aspmx.l.google.com."),
  MX("@", 5, "alt1.aspmx.l.google.com."),
  MX("@", 5, "alt2.aspmx.l.google.com."),
  MX("@", 10, "alt3.aspmx.l.google.com."),
  MX("@", 10, "alt4.aspmx.l.google.com."),
];

// Every zone uses Cloudflare's automatic TTL, which its API writes as 1.
DEFAULTS(DefaultTTL(1));

D("poly.red", REG_NONE, DnsProvider(DSP_CLOUDFLARE),
  A("@", SERVER, CF_PROXY_ON),
  GOOGLE_MAIL,
  TXT("@", "v=spf1 include:_spf.google.com ~all"),
  TXT("@", "google-site-verification=_f_PMaAXzRTN7CaQ82cUuz-6VAfS5qfWNVpKqgkyhTE"),
  TXT("google._domainkey", "v=DKIM1; k=rsa; p=MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAq51RstHyjefzzeopgLOnwXzreFqH3X836pINDbtsVdKAHXXdoTVvSdACBlLrA3TqxK1a/nvfvc9EcFBiLE3hjUao66eA/HswMeR0HrdKvQVG5xnyAkZg7o4MyG4PzcloragP1F6GVliR2Zb7rqDRjaudZL8Aq7uFMk8pDVPZwvG14qwtoom+O9bdio1qWtJqQyOUWFz1OHcTzpbFNoZLj0HXkvNB4gkMCu2rT5wu8ATEvszP6+HwqZw70XaG0Y2/zI7nplFqNBRvXvSvVI3KFGlB2IuvI8SDbvdFHYfVV5CtrQgDhgiwRTGXohQLxReC3kUNQEAajhXzDLfA1rx0PQIDAQAB"),
);

// changkun.us only redirects to the blog, path and query kept. Cloudflare
// answers the redirect, so the hosts are proxied but nothing reaches the
// server.
D("changkun.us", REG_NONE, DnsProvider(DSP_CLOUDFLARE),
  A("@", SERVER, CF_PROXY_ON),
  A("ou", SERVER, CF_PROXY_ON),
  CNAME("www", "changkun.us.", CF_PROXY_ON),
  CF_SINGLE_REDIRECT("changkun.us to the blog", 301,
    'http.host in {"changkun.us" "www.changkun.us" "ou.changkun.us"}',
    'concat("https://changkun.de/blog", http.request.uri.path)'),
  GOOGLE_MAIL,
  TXT("@", "v=spf1 include:_spf.google.com ~all"),
  TXT("@", "google-site-verification=1Q8wu_zPZ1V11olfHzIHAmfzTe-WxOGXY5DDox5qjzE"),
  TXT("google._domainkey", "v=DKIM1; k=rsa; p=MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAvF2JohXnnCCMhM2/083CLWQ5rAAisiQSEjHC0luCMMUSWghFE2dmbWWEts3N+FPgop0kCYo33iRmBNRrxRMA//Iz+EN0Un0PnUisKAd3xW45chviiEJFMFDz0bNFoa6NEYw7FJpwxj6uhRa1OccbLsBB5XyrKzmTYhUo5Pk37ABH9cU0JpLCLnJbJ/9UIq7uREWNG82Xag268wVDXHwne2VwddbZLmjuktQMHvcmeJ2PoRkoPyptrVBVLTMd5hi2juM1pGrL7+iteYVbmabGCQj8Psa8eHavQgneLMxsf4epbIB/7YUmRKMnEyA4tuVJaQsoJBfWzYC2sYCAvadkmwIDAQAB"),
);

D("golang.design", REG_NONE, DnsProvider(DSP_CLOUDFLARE),
  A("@", SERVER, CF_PROXY_ON),
  A("www", SERVER, CF_PROXY_ON),
  GOOGLE_MAIL,
  TXT("@", "v=spf1 include:_spf.google.com ~all"),
  TXT("@", "google-site-verification=OwmHdWz2c67bUMSMcA5uobnyoKe0Z3-GLLPC9e0WwbI"),
  TXT("@", "D8008724"),
  TXT("google._domainkey", "v=DKIM1; k=rsa; p=MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAothqatbrOzTPVC1/cfNFyLpvoawRU1u7wzk+YYUoFL7xyOIAOSs43z50pA3sSSDPhTTjLKVcaF35BsBTWCDTrDyiJxu8l4pCTi7cJGd9dhyDYY/JSN1/+ABPeEL7CGbovGbRbPglIC1jOdqZlArgK9OYlBEkxuQwdZ4qE8sroasx8RLRimNNZZEF33f16m4m0B1XYcC0AWvcSbtsl1FsFYrACRZhLCz4SfW+W7IeRNKdJb2NOgMZqTCyC3IrYgU/b8GJF2lB00n14U4uZSgZZGFE78cg8sCafPrq0G/dRannRrYkyaoyhmIQ9TnQFEeJr/e/jCjSG/1BBDicIlg1dwIDAQAB"),
  // DMARC reports go to Cloudflare's DMARC Management.
  TXT("_dmarc", "v=DMARC1; p=none; rua=mailto:a1c2e672918a4351b65309326ce13e73@dmarc-reports.cloudflare.net"),
);

// changkun.de is not proxied: ssh and the publishing workflows reach the
// server by this name, and Cloudflare's proxy does not carry ssh.
D("changkun.de", REG_NONE, DnsProvider(DSP_CLOUDFLARE),
  A("@", SERVER),
  A("www", SERVER),
  A("ou", SERVER),
  A("api", SERVER),
  A("blog", SERVER),
  // lmu and mail only redirect, and Cloudflare answers them, so their
  // address is a placeholder nothing reaches.
  A("lmu", "192.0.2.1", CF_PROXY_ON),
  A("mail", "192.0.2.1", CF_PROXY_ON),
  CF_SINGLE_REDIRECT("lmu to the LMU team page", 302,
    'http.host eq "lmu.changkun.de"',
    'concat("https://www.medien.ifi.lmu.de/team/changkun.ou/", "")'),
  CF_SINGLE_REDIRECT("mail to Gmail", 302,
    'http.host eq "mail.changkun.de"',
    'concat("https://mail.google.com/", "")'),
  GOOGLE_MAIL,
  TXT("@", "v=spf1 include:_spf.google.com ~all"),
  TXT("@", "google-site-verification=cAfLSlwlWzoE5qEJv0U6ExGy8Ytu-vUihx_i9g1uPuY"),
  TXT("@", "openai-domain-verification=dv-hBYFVae7AHo9eiRl7zLdndFD"),
  TXT("@", "D9731473"),
  TXT("google._domainkey", "v=DKIM1; k=rsa; p=MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAg7at5QBS/3PYYfqqkbpxEJfIh/x4s6a6ojw/MGgRwJ5356MOX/0hAi7lv2GEzYVkLeEZwjTmBT9sDgnjCClnOm5tg7VnZBGiOMRTkyNwHKf+TS8RWVn0FzHwiLXSAN1mAx+YxbaHvFX0JoF1IgSLJpB0xzucUTmLT37z7DBv4LFUJfsR6YZ5O0cYu/8bfImC/mtOxN7tdNd3damI7379EMVZt03C2DHUVZ+yfSg1HEJh8of/yTKclsXOWsoU49oE7Cd8lOhP7CYbuVzz+qe3lcarylqH/kkFnH6GqJBLnW1dt8h3/ZZ8C7sfc47B1UpHYmAWeMvW5dEPKRN6Yb4W6wIDAQAB"),
  TXT("_dmarc", "v=DMARC1; p=none; rua=mailto:dmarc-reports@changkun.de"),
  // Mailgun sends as mailgun.changkun.de.
  MX("mailgun", 10, "mxa.eu.mailgun.org."),
  MX("mailgun", 10, "mxb.eu.mailgun.org."),
  TXT("mailgun", "v=spf1 include:eu.mailgun.org ~all"),
  TXT("mta._domainkey.mailgun", "k=rsa; p=MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDA0yVxCr0+nb9b1/ScsJOxQnpoi8FzmhEdnpReWl/UjzxEx/NOsPPUhVjzp2QU1MwZb8Cs8sC87n50fSU/Q0fYUU8hDGMk3VlP+snAC+wGAFx+fGDxJ7B+O8SFKwzJIgs+x0trlTrBOJK+5eGmtc+i6opyNuYY8dyhruZYuvKYsQIDAQAB"),
  CNAME("email.mailgun", "eu.mailgun.org."),
);
