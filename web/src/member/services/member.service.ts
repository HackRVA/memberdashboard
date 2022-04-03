// rxjs
import { Observable } from 'rxjs';

// memberdashboard
import { MemberResponse } from '../types/api/member-response';
import { CreateMemberRequest } from '../types/api/create-member-request';
import { UpdateMemberRequest } from './../types/api/update-member-request';
import { AckResponse } from './../../shared/types/api/ack-response';
import { AssignRFIDRequest } from '../types/api/assign-rfid-request';
import { HTTPService } from '../../shared/services/http.service';
import { ENV } from '../../env';
import { Injectable } from '../../shared/di';

@Injectable('member')
export class MemberService extends HTTPService {
  private readonly memberUrlSegment: string = ENV.api + '/member';

  getMembers(): Observable<MemberResponse[]> {
    return this.get<MemberResponse[]>(this.memberUrlSegment);
  }

  updateMemberByEmail(
    email: string,
    request: UpdateMemberRequest
  ): Observable<AckResponse> {
    return this.put<AckResponse>(
      `${this.memberUrlSegment}/email/${email}`,
      request
    );
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