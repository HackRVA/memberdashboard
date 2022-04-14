// memberdashboard
import { MemberLevel } from './types/custom/member-level';
import { Member } from './types/api/member-response';

import structuredClone from '@ungap/structured-clone'; // polyfill


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

function getResourcesLabel(member: Member) {
  if (!member.resources) return member;

  // return structuredClone({
  //   resourcesLabel: member.resources.map(resource => resource.name).join(','),
  //   resources: structuredClone([...member.resources]),
  //   ...member,
  // });
  return {
    resourcesLabel: member.resources.map(resource => resource.name).join(','),
    resources: [...member.resources],
    ...member,
  };
}

function withResourceLabels(members: Member[]): Member[] {
  return members.map(getResourcesLabel)
}

export function deepCopy(members: Member[]) {
  return structuredClone(withResourceLabels(members));;
}
