import { AbstractControl, FormGroup, ValidationErrors } from '@angular/forms';

export function passwordMatchValidator(
  registerFormGroup: FormGroup
): ValidationErrors | null {
  const passwordCtrl: AbstractControl = registerFormGroup.get('password');
  const retypePasswordCtrl: AbstractControl =
    registerFormGroup.get('retypePassword');

  if (!passwordCtrl.value || !retypePasswordCtrl.value) {
    return null;
  }

  const password: string = passwordCtrl.value.trim();
  const retypePassword: string = retypePasswordCtrl.value.trim();

  if (password.length === 0 || password.length === 0) {
    return null;
  }

  return password === retypePassword
    ? null
    : ({ mismatch: true } as ValidationErrors);
}
