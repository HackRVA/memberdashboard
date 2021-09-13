// rxjs
import { Observable } from 'rxjs';

// memberdashboard
import { MemberResponse } from '../types/api/member-response';
import { CreateMemberRequest } from '../types/api/create-member-request';
import { AssignRFIDRequest } from '../types/api/assign-rfid-request';
import { HTTPService } from '../../shared/services/http.service';
import { ENV } from '../../env';

export class MemberService extends HTTPService {
  private readonly memberUrlSegment: string = ENV.api + '/member';

  getMembers(): Observable<MemberResponse[]> {
    return this.get<MemberResponse[]>(this.memberUrlSegment);
  }

  assignRFID(request: AssignRFIDRequest): Observable<void> {
    return this.post<void>(this.memberUrlSegment + '/assignRFID', request);
  }

  assignNewMemberRFID(request: CreateMemberRequest): Observable<void> {
    return this.post<void>(this.memberUrlSegment + '/new', request);
  }

  assignRFIDToSelf(request: AssignRFIDRequest): Observable<void> {
    return this.post<void>(this.memberUrlSegment + '/assignRFID/self', request);
  }

  getMemberByEmail(email: string): Observable<MemberResponse> {
    return this.get<MemberResponse>(`${this.memberUrlSegment}/email/${email}`);
  }

  getUsersMemberInfo(): Observable<MemberResponse> {
    return this.get<MemberResponse>(`${this.memberUrlSegment}/self`);
  }

  downloadNonMembersCSV(): Observable<Blob> {
    return this.get<Blob>(`${this.memberUrlSegment}/slack/nonmembers`, 'blob');
  }
}
