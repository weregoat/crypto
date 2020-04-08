#ECB cut-and-paste

Write a k=v parsing routine, as if for a structured cookie. The routine should take:

```
foo=bar&baz=qux&zap=zazzle
```
... and produce:

```
{
  foo: 'bar',
  baz: 'qux',
  zap: 'zazzle'
}
```
(you know, the object; I don't care if you convert it to JSON).

Now write a function that encodes a user profile in that format, given an email address. You should have something like:
```
profile_for("foo@bar.com")
```
... and it should produce:
```
{
  email: 'foo@bar.com',
  uid: 10,
  role: 'user'
}
```
... encoded as:
```
email=foo@bar.com&uid=10&role=user
```
Your "profile_for" function should not allow encoding metacharacters (& and =). Eat them, quote them, whatever you want to do, but don't let people set their email address to "foo@bar.com&role=admin".

Now, two more easy functions. Generate a random AES key, then:

* Encrypt the encoded user profile under the key; "provide" that to the "attacker".
* Decrypt the encoded user profile and parse it.

Using only the user input to profile_for() (as an oracle to generate "valid" ciphertexts) and the ciphertexts themselves, make a role=admin profile.

## Goat notes

On first take it seems a way to solve this is to submit one (or plausibly more) email address so as to get back the code blocks for (or, more properly, containing) "admin", replace it, and send it back.

It seems a very specific example, and I can imagine a scenario for it. Which makes me think:
* What I wrote earlier: attack targets specific implementations.
* How on Earth could anyone think that storing such information in a cookie was a sane and secure idea?

On the other hand, what we have here is basically a token.

Another thought: this attack assume we know a lot about the encrypted string. We know, it has `role` and that we want `admin`, that is at the end...

I suppose that you can get to this kind of information in many ways, but it also occurred to me that, possibly, you get them by trying a lot of plausible combinations.

Which reminds me that article I read a long time ago about that one of the best defence against password attacks was to limit the number of allowed tries (or slow them down a lot); I think it was from Microsoft people.
