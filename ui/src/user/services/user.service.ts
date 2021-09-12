import { UserResponse } from './../types/api/user-response';
// rxjs
import { Observable } from 'rxjs';
import { ENV } from '../../env';
import { HTTPService } from '../../shared/services/http.service';

// memberdashboard

export class UserService extends HTTPService {
  private readonly userUrlSegment: string = ENV.api + '/user';

  getUser(): Observable<UserResponse> {
    return this.get<UserResponse>(this.userUrlSegment);
  }
}
