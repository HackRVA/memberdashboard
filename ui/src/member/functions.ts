import { MemberLevel } from './types/custom/member-level';

export const displayMemberStatus = (memberLevel: MemberLevel): string => {
  switch (memberLevel) {
    case MemberLevel.inactive:
      return 'Inactive';
    case MemberLevel.credited:
      return 'Credited';
    case MemberLevel.classic:
      return 'Classic';
    case MemberLevel.standard:
      return 'Standard';
    case MemberLevel.premium:
      return 'Premium';
    default:
      return 'No member status found';
  }
};
