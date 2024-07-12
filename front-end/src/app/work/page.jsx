// app/work/page.jsx
"use client";
import { useState, useEffect, useCallback } from "react";
import Link from "next/link"; // Import Link component for navigation

export default function WorkPage() {
  const [tasks, setTasks] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState(null);

  const fetchTasks = useCallback(async () => {
    try {
      const token = localStorage.getItem("token");
      const response = await fetch("http://localhost:8080/api/work/today", {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      if (!response.ok) {
        throw new Error(`Request failed with status ${response.status}`);
      }

      const data = await response.json();
      setTasks(data);
    } catch (error) {
      setError(error.message);
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchTasks();
  }, [fetchTasks]);

  const handleToggleComplete = async (taskId) => {
    try {
      const token = localStorage.getItem("token");
      const updatedTask = tasks.find((t) => t.id === taskId);

      const response = await fetch(`http://localhost:8080/api/work/${taskId}`, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          taskId: taskId,
          completed: !updatedTask.daily_completed,
          date: new Date().toISOString().slice(0, 10),
        }),
      });

      if (!response.ok) {
        throw new Error(`Request failed with status ${response.status}`);
      }

      fetchTasks();
    } catch (error) {
      setError(error.message);
    }
  };

  const today = new Date().toLocaleDateString();

  return (
    <div className="container mx-auto mt-8">
      WorkHistoryPage
      <h2 className="text-2xl font-bold mb-4">Work Log for {today}</h2>
      {isLoading ? (
        <p>Loading tasks...</p>
      ) : error ? (
        <p className="text-red-500">Error: {error}</p>
      ) : (
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              <th
                scope="col"
                className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
              >
                No.
              </th>
              <th
                scope="col"
                className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
              >
                Title
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
            {tasks.map((task, index) => (
              <tr key={task.id}>
                <td className="px-6 py-4 whitespace-nowrap">{index + 1}</td>
                <td className="px-6 py-4 whitespace-nowrap">{task.title}</td>
                <td className="px-6 py-4 whitespace-nowrap">
                  {task.description}
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <input
                    type="checkbox"
                    checked={task.daily_completed}
                    onChange={() => handleToggleComplete(task.id)}
                  />
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
}
