// rxjs
import { Observable } from 'rxjs';

// memberdashboard
import { UserResponse } from './../types/api/user-response';
import { ENV } from '../../env';
import { HTTPService } from '../../shared/services/http.service';
import { Injectable } from '../../shared/di';

@Injectable('user')
export class UserService extends HTTPService {
  private readonly userUrlSegment: string = ENV.api + '/user';

  getUser(): Observable<UserResponse> {
    return this.get<UserResponse>(this.userUrlSegment);
  }
}
