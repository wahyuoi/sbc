
#todo
- [ ] Add error for conflict email
- [ ] Add error for invalid email
- [ ] Add error for invalid audio format (limit to wav & m4a)
- [ ] Add error for invalid audio file duration
- [ ] Add error for invalid phrase ID when submitting audio
- [ ] Add error for invalid user ID when submitting audio
- [ ] Handle resubmit audio

Notes:
- // RemoveAll removes any temporary files associated with a Form.
func (f *Form) RemoveAll() error {
