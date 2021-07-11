// rxjs
import { BehaviorSubject } from "rxjs";

// memberdashboard
import { AuthUserProfile } from "./components/shared/types";

export const authUser = new BehaviorSubject<AuthUserProfile>({
  login: false,
  email: null,
});
