import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import {
  AckResponse,
  AssignRFIDRequest,
  CreateMemberRequest,
  MemberResponse,
  MemberSearchRequest,
  UpdateMemberRequest,
} from '../types';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class MemberService {
  private readonly _memberUrlSegment: string = '/api/member';
  constructor(private readonly http: HttpClient) {}

  getMembers(request: MemberSearchRequest): Observable<MemberResponse[]> {
    const params: URLSearchParams = new URLSearchParams();

    if (typeof request.page === 'number') {
      params.append('page', request.page.toString());
    }

    if (request.count > 0) {
      params.append('count', request.count.toString());
    }

    if (typeof request.active === 'boolean') {
      params.append('active', request.active.toString());
    }

    if (request.search) {
      params.append('search', request.search);
    }

    return this.http.get<MemberResponse[]>(
      `${this._memberUrlSegment}?${params.toString()}`
    );
  }

  getMemberSelf(): Observable<MemberResponse> {
    return this.http.get<MemberResponse>(this._memberUrlSegment + '/self');
  }

  updateMemberByEmail(
    email: string,
    request: UpdateMemberRequest
  ): Observable<AckResponse> {
    return this.http.put<AckResponse>(
      `${this._memberUrlSegment}/email/${email}`,
      request
    );
  }

  assignRFID(request: AssignRFIDRequest): Observable<void> {
    return this.http.post<void>(
      this._memberUrlSegment + '/assignRFID',
      request
    );
  }

  assignRFIDToSelf(request: AssignRFIDRequest): Observable<void> {
    return this.http.post<void>(
      this._memberUrlSegment + '/assignRFID/self',
      request
    );
  }

  assignNewMemberRFID(request: CreateMemberRequest): Observable<void> {
    return this.http.post<void>(this._memberUrlSegment + '/new', request);
  }

  checkMemberStatus(subscriptionID: string): Observable<MemberResponse> {
    return this.http.get<MemberResponse>(
      `${this._memberUrlSegment}/${subscriptionID}/status`
    );
  }
}
