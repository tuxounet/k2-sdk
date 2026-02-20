import React, { useState, useEffect } from "react";
import DigitalClock from "../components/DigitalClock";
import { getTime } from "../api/time";
import { TimeResponse } from "../types/time";
import "./ClockPage.css";

const ClockPage: React.FC = () => {
  const [time, setTime] = useState<TimeResponse | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  const fetchTime = async () => {
    try {
      const data = await getTime();
      setTime(data);
      setError(null);
      setLoading(false);
    } catch (err) {
      setError("Erreur lors de la récupération de l'heure");
      setLoading(false);
      console.error("Error fetching time:", err);
    }
  };

  useEffect(() => {
    // Sync with server every minute
    const syncInterval = setInterval(fetchTime, 1000);

    return () => {
      clearInterval(syncInterval);
    };
  }, [time]);

  if (loading) {
    return (
      <div className="clock-page">
        <div className="loading-message">Chargement...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="clock-page">
        <div className="error-message">{error}</div>
      </div>
    );
  }

  if (!time) {
    return null;
  }

  return (
    <div className="clock-page">
      <div className="clock-container">
        <div className="clock-display">
          <DigitalClock
            hours={time.hours}
            minutes={time.minutes}
            seconds={time.seconds}
            day={time.day}
            month={time.month}
            year={time.year}
            dayOfWeek={time.dayOfWeek}
          />
        </div>
      </div>
    </div>
  );
};

export default ClockPage;
