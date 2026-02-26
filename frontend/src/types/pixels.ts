export interface RetargetingPixel {
  id: string;
  link_id: string;
  pixel_type: PixelType;
  pixel_id?: string;
  custom_script?: string;
  is_active: boolean;
  created_at: string;
}

export type PixelType = 'facebook' | 'google_ads' | 'tiktok' | 'linkedin' | 'custom';

export interface CreatePixelRequest {
  pixel_type: PixelType;
  pixel_id?: string;
  custom_script?: string;
}

export const PIXEL_TYPES = [
  {
    value: 'facebook' as PixelType,
    label: 'Meta (Facebook) Pixel',
    icon: '📘',
    placeholder: 'e.g. 123456789012345',
    description: 'Tracks visitors for Meta ads retargeting',
    requiresId: true,
  },
  {
    value: 'google_ads' as PixelType,
    label: 'Google Ads',
    icon: '🔍',
    placeholder: 'e.g. AW-123456789',
    description: 'Google Ads global site tag for remarketing',
    requiresId: true,
  },
  {
    value: 'tiktok' as PixelType,
    label: 'TikTok Pixel',
    icon: '🎵',
    placeholder: 'e.g. ABCDE1234567890',
    description: 'TikTok Pixel for audience building',
    requiresId: true,
  },
  {
    value: 'linkedin' as PixelType,
    label: 'LinkedIn Insight Tag',
    icon: '💼',
    placeholder: 'e.g. 1234567',
    description: 'LinkedIn Insight Tag for B2B retargeting',
    requiresId: true,
  },
  {
    value: 'custom' as PixelType,
    label: 'Custom Script',
    icon: '⚙️',
    placeholder: '<script>/* your tracking code */</script>',
    description: 'Inject arbitrary tracking code (use with care)',
    requiresId: false,
  },
] as const;
