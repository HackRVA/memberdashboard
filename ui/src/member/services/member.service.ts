// rxjs
import { Observable } from 'rxjs';
import { ENV } from '../../env';
import { HTTPService } from '../../shared/services/http.service';
import { AssignRFIDRequest } from '../types/api/assign-rfid-request';
import { CreateMemberRequest } from '../types/api/create-member-request';
import { MemberResponse } from '../types/api/member-response';

// memberdashboard

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
    return this.get<Blob>(`${this.memberUrlSegment}/slack/nonmembers`);
  }
}
