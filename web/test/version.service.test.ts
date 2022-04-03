// testing
import { expect, assert } from '@open-wc/testing';

// rxjs
import { of } from 'rxjs';

// memberdashboard
import { VersionResponse } from './../src/shared/types/api/version-response';
import { VersionService } from './../src/shared/services/version.service';

describe('VersionService', () => {
  let mockService: VersionService;

  before(async () => {
    mockService = setupMockService();
  });

  it('should be defined', () => {
    assert.instanceOf(mockService, VersionService);
  });

  it('should get version number of the app', () => {
    // ARRANGE
    let actualResponse: VersionResponse = null;

    // ACT
    mockService.getVersion().subscribe((response: VersionResponse) => {
      actualResponse = response;
    });

    // ASSERT
    expect(actualResponse).not.be.null;
    expect(actualResponse.major).equal('1');
    expect(actualResponse.build).equal('abc');
  });
});

function setupMockService(): VersionService {
  let service = new VersionService();

  service.getVersion = () => {
    return of({
      major: '1',
      build: 'abc',
    });
  };

  return service;
}
