// app/work-history/page.jsx
"use client";
import { useState, useEffect, useCallback } from "react";

export default function WorkHistoryPage() {
  const [workLogs, setWorkLogs] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState(null);

  const fetchWorkLogHistory = useCallback(async () => {
    try {
      const token = localStorage.getItem("token");
      const response = await fetch("http://localhost:8080/api/work", {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      if (!response.ok) {
        throw new Error(`Request failed with status ${response.status}`);
      }

      const data = await response.json();
      setWorkLogs(data);
    } catch (error) {
      setError("Failed to fetch work log history. Please try again.");
    } finally {
      setIsLoading(false);
    }
  }, []);

  const groupedWorkLogs = workLogs.reduce((acc, log) => {
    const date = new Date(log.date).toLocaleDateString();
    acc[date] = acc[date] || [];
    acc[date].push(log);
    return acc;
  }, {});

  const [expandedDates, setExpandedDates] = useState({});

  const toggleDate = (date) => {
    setExpandedDates((prevExpanded) => ({
      ...prevExpanded,
      [date]: !prevExpanded[date],
    }));
  };

  useEffect(() => {
    fetchWorkLogHistory();
  }, [fetchWorkLogHistory]);

  return (
    <div className="container mx-auto mt-8">
      <h2 className="text-2xl font-bold mb-4">Work Log History</h2>

      {/* Error Handling and Loading State */}
      {error ? (
        <p className="text-red-500">Error: {error}</p>
      ) : isLoading ? (
        <p>Loading history...</p>
      ) : (
        <div>
          {/* Conditional Check for Data */}
          {Object.keys(groupedWorkLogs).length >0 ? (
            // Grouped work logs
            Object.entries(groupedWorkLogs).map(([date, logs]) => (
              <div key={date}>
                <button
                  onClick={() => toggleDate(date)}
                  className="w-full flex justify-between items-center px-4 py-2 text-left text-gray-500 bg-gray-100 rounded-t-lg focus:outline-none"
                >
                  <span className="font-medium">{date}</span>
                  {/* Toggle Arrow */}
                  <svg
                    className={`w-4 h-4 transform ${expandedDates[date] ? 'rotate-180' : ''}`}
                    xmlns="http://www.w3.org/2000/svg"
                    viewBox="0 0 20 20"
                    fill="currentColor"
                    aria-hidden="true"
                  >
                    <path
                      fillRule="evenodd"
                      d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z"
                      clipRule="evenodd"
                    />
                  </svg>
                </button>
                {expandedDates[date] && (
                  <table className="min-w-full divide-y divide-gray-200">
                    {/* Table Headers */}
                    <thead>
                      <tr>
                        <th
                          scope="col"
                          className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                        >
                          Task
                        </th>
                        <th
                          scope="col"
                          className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                        >
                          Description
                        </th>
                        <th
                          scope="col"
                          className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                        >
                          Completed
                        </th>
                      </tr>
                    </thead>
                    <tbody className="bg-white divide-y divide-gray-200">
                      {logs.map((log) => (
                        <tr key={log.id}>
                          <td className="px-6 py-4 whitespace-nowrap">
                            {log.title}
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap">
                            {log.description}
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap">
                            {log.completed ? "Yes" : "No"}
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                )}
              </div>
            ))
          ) : (
            <p className="text-gray-500 mt-4">No work log history found.</p>
          )}
        </div>
      )}
    </div>
  );
}
