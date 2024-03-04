import { Pipe, PipeTransform } from '@angular/core';
import { MemberLevel } from '../types';

@Pipe({
  name: 'memberLevel',
  standalone: true,
})
export class MemberLevelPipe implements PipeTransform {
  transform(value: MemberLevel): string {
    switch (value) {
      case MemberLevel.Inactive:
        return MemberLevel[MemberLevel.Inactive];
      case MemberLevel.Credited:
        return MemberLevel[MemberLevel.Credited];
      case MemberLevel.Classic:
        return MemberLevel[MemberLevel.Classic];
      case MemberLevel.Standard:
        return MemberLevel[MemberLevel.Standard];
      case MemberLevel.Premium:
        return MemberLevel[MemberLevel.Premium];
      default:
        return 'No member status found';
    }
  }
}
