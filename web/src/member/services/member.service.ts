// rxjs
import { Observable } from 'rxjs';

// memberdashboard
import { Member } from '../types/api/member-response';
import { CreateMemberRequest } from '../types/api/create-member-request';
import { UpdateMemberRequest } from './../types/api/update-member-request';
import { AckResponse } from './../../shared/types/api/ack-response';
import { AssignRFIDRequest } from '../types/api/assign-rfid-request';
import { HTTPService } from '../../shared/services/http.service';
import { ENV } from '../../env';
import { Injectable } from '../../shared/di';
import { Inject } from '../../shared/di/inject';
import { deepCopy } from '../functions';

@Injectable('member')
export class MemberService extends HTTPService {
  private readonly memberUrlSegment: string = ENV.api + '/member';

  getMembers(): Observable<Member[]> {
    return this.get<Member[]>(this.memberUrlSegment);
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

  getMemberByEmail(email: string): Observable<Member> {
    return this.get<Member>(`${this.memberUrlSegment}/email/${email}`);
  }

  getUsersMemberInfo(): Observable<Member> {
    return this.get<Member>(`${this.memberUrlSegment}/self`);
  }

  downloadNonMembersCSV(): Observable<Blob> {
    return this.get<Blob>(`${this.memberUrlSegment}/slack/nonmembers`, 'blob');
  }
}

@Injectable('member-manager')
export class MemberManagerService {
  @Inject('member')
  private memberService: MemberService;

  allMembers: Member[];
  inactiveMembers: Member[];
  activeMembers: Member[];
  _showActive: boolean = true;

  _filteredMembers: Member[] = [];
  _refreshMembers: () => void;

  get filteredMembers(): Member[] {
    return this._filteredMembers;
  }
  set filteredMembers(members: Member[]) {
    this._filteredMembers = deepCopy(members);
    this.updateListeners();
  }

  get showActive(): boolean {
    return this._showActive;
  }
  set showActive(showActive: boolean) {
    this._showActive = showActive;
    if (!showActive) {
      this.filteredMembers = this.inactiveMembers;
      return;
    }
    this.filteredMembers = this.activeMembers;
  }

  listeners: Array<() => void> = [];

  registerListener = (listener: () => void): void => {
    this.listeners.push(listener);
  };

  updateListeners(): void {
    this.listeners.forEach(listener => listener());
  }

  updateMembers(memberResponse: Member[]): void {
    this.allMembers = memberResponse;
    this.inactiveMembers = this.allMembers.filter(
      member => member.memberLevel === 1
    );
    this.activeMembers = this.allMembers.filter(
      member => member.memberLevel > 1
    );

    const members = {
      all: this.allMembers,
      active: this.activeMembers,
      inactive: this.inactiveMembers,
    };

    this.filteredMembers = this.showActive ? members.active : members.inactive;
  }

  getMembers = (): Member[] => {
    this.memberService.getMembers().subscribe(
      (members: Member[]) => {
        this.filteredMembers = members;
        console.log(members.length);
      },
      (error: null) => {
        console.log(error);
      }
    );
    this.updateMembers(this.filteredMembers);
    return this.filteredMembers;
  };
}

export default new MemberManagerService();
