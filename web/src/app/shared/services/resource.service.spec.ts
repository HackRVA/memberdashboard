import { TestBed } from '@angular/core/testing';
import { HttpClient } from '@angular/common/http';
import { of } from 'rxjs';
import { ResourceService } from './resource.service';
import { SharedSpies } from '../testings';

describe('ResourceService', () => {
  let service: ResourceService;
  let http: jasmine.SpyObj<HttpClient>;
  const resourceUrlSegment: string = '/api/resource';

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [
        { provide: HttpClient, useValue: SharedSpies.createHttpClientSpy() },
      ],
    });
    service = TestBed.inject(ResourceService);
    http = TestBed.inject(HttpClient) as jasmine.SpyObj<HttpClient>;
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  it('should update ACLs', () => {
    // ARRANGE
    http.post.and.returnValue(of(null));

    // ACT
    service.updateACLs().subscribe();

    // ASSERT
    expect(http.post).toHaveBeenCalledWith(
      `${resourceUrlSegment}/updateacls`,
      {}
    );
  });

  it('should remove ACLs', () => {
    // ARRANGE
    http.delete.and.returnValue(of(null));

    // ACT
    service.removeACLs().subscribe();

    // ASSERT
    expect(http.delete).toHaveBeenCalledWith(
      `${resourceUrlSegment}/deleteacls`
    );
  });
});
