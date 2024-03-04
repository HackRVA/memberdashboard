import { TestBed } from '@angular/core/testing';
import { VersionService } from './version.service';
import { HttpClient } from '@angular/common/http';
import { SharedSpies } from '../testings';
import { VersionResponse } from '../types';
import { of } from 'rxjs';

describe('VersionService', () => {
  let service: VersionService;
  let http: jasmine.SpyObj<HttpClient>;
  const versionUrlSegment: string = '/api/version';

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [
        { provide: HttpClient, useValue: SharedSpies.createHttpClientSpy() },
      ],
    });
    service = TestBed.inject(VersionService);
    http = TestBed.inject(HttpClient) as jasmine.SpyObj<HttpClient>;
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  it('should get version', () => {
    // ARRANGE
    let mockResponse: VersionResponse = {
      major: '1',
      build: 'dev',
    };
    let expectedResponse: VersionResponse = null;
    http.get.and.returnValue(of(mockResponse));

    // ACT
    service.getVersion().subscribe((response: VersionResponse) => {
      expectedResponse = response;
    });

    // ASSERT
    expect(expectedResponse).not.toBe(null);
    expect(expectedResponse.major).toBe(mockResponse.major);
    expect(expectedResponse.build).toBe(mockResponse.build);
    expect(http.get).toHaveBeenCalledWith(versionUrlSegment);
  });
});
