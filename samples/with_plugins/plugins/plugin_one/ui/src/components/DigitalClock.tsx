import React from "react";
import "./DigitalClock.css";

interface DigitalClockProps {
  hours: number;
  minutes: number;
  seconds: number;
  day: number;
  month: number;
  year: number;
  dayOfWeek: string;
}

const DigitalClock: React.FC<DigitalClockProps> = ({
  hours,
  minutes,
  seconds,
  day,
  month,
  year,
  dayOfWeek,
}) => {
  const formatNumber = (num: number): string => {
    if (!num && num !== 0) {
      return "--";
    }
    
    return num.toString().padStart(2, "0");
  };

  return (
    <div className="digital-clock">
      <div className="digital-time">
        {formatNumber(hours)}:{formatNumber(minutes)}:{formatNumber(seconds)}
      </div>
      <div className="digital-date">
        {formatNumber(day)}/{formatNumber(month)}/{year}
      </div>
      <div className="digital-day">{dayOfWeek}</div>
    </div>
  );
};

export default DigitalClock;
