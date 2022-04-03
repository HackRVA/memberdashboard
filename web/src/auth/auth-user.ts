// rxjs
import { BehaviorSubject } from 'rxjs';

// memberdashboard
import { AuthUserProfile } from './types/custom/auth-user-profile';

export const authUser$ = new BehaviorSubject<AuthUserProfile>({
  login: false,
  email: null,
});
