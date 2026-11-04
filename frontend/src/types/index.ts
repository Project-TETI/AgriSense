export interface SummaryMetric {
  id: string;
  title: string;
  value: string;
  changeText: string;
  changeType: 'up' | 'down';
  icon: 'search' | 'shield' | 'alert';
}

export interface ActivityDay {
  day: string;
  count: number;
  heightPx: number;
  isPeak?: boolean;
}

export interface DiseaseEducationItem {
  id: string;
  title: string;
  description: string;
  imageUrl: string;
  actionText: string;
}

export interface InspectionRecord {
  id: string;
  cropName: string;
  conditionDetail: string;
  dateStr: string;
  statusBadge: {
    label: string;
    variant: 'warning' | 'success';
  };
  accuracyPercent: number;
  actionButton: {
    label: string;
    variant: 'action' | 'optimal';
  };
  imageUrl: string;
}
