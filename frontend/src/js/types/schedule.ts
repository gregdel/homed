// Schedule types

export interface TimeSlot {
  id: string;
  start: string; // "HH:MM" format
  value: string | number;
}

export interface DailySchedule {
  weekday: number; // 0-6 (Sunday-Saturday)
  time_slots: TimeSlot[];
}

export interface ScheduleOverride {
  id: string;
  date: string; // ISO date
  value: string | number;
}

export interface Schedule {
  default_value: string | number;
  overrides: ScheduleOverride[];
  daily: DailySchedule[];
}
