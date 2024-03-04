import { TestBed } from '@angular/core/testing';
import { MemberService } from './member.service';
import { HttpClient } from '@angular/common/http';
import { SharedSpies } from '../testings';
import { MemberResponse } from '../types';
import { of } from 'rxjs';

describe('MemberService', () => {
  let service: MemberService;
  let http: jasmine.SpyObj<HttpClient>;

  const memberUrlSegment: string = '/api/member';

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [
        { provide: HttpClient, useValue: SharedSpies.createHttpClientSpy() },
      ],
    });
    service = TestBed.inject(MemberService);
    http = TestBed.inject(HttpClient) as jasmine.SpyObj<HttpClient>;
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  it('should get self', () => {
    // ARRANGE
    const id: string = '123';
    let expectedResponse: MemberResponse = null;
    http.get.and.returnValue(of({ id: id } as MemberResponse));

    // ACT
    service.getMemberSelf().subscribe((response: MemberResponse) => {
      expectedResponse = response;
    });

    // ASSERT
    expect(expectedResponse).not.toBe(null);
    expect(expectedResponse.id).toBe(id);
    expect(http.get).toHaveBeenCalledWith(`${memberUrlSegment}/self`);
  });
});
