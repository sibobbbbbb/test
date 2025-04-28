import { prisma } from '../config/prisma.js';
import { successResponse, errorResponse } from '../utils/response.js';
import logger from '../utils/logger.js';

export const getCommunityEvents = async (req, res) => {
  try {
    let events;
    // Jika user memiliki role "Buddy", hanya tampilkan event dengan accessLevel "MemberAndBuddy"
    if (req.user && req.user.access === 'Buddy') {
      events = await prisma.event.findMany({
        where: {
          accessLevel: 'MemberAndBuddy',
          eventType: 'Community',
        },
      });
    } else {
      events = await prisma.event.findMany({
        where: {
          eventType: 'Community',
        },
      });
    }
    return successResponse(res, 200, 'Success', events);
  } catch (error) {
    return errorResponse(res, 500, error.message);
  }
};

export const getCommunityEvent = async (req, res) => {
  try {
    const event = await prisma.event.findUnique({
      where: {
        id: parseInt(req.params.id),
      },
      include: {
        attendees: {
          include: {
            user: {
              select: {
                name: true,
                email: true,
              },
            },
          },
        },
      },
    });
    if (!event || event.eventType !== 'Community') {
      return errorResponse(res, 404, 'Community Event not found');
    }
    return successResponse(res, 200, 'Success', event);
  } catch (error) {
    return errorResponse(res, 500, error.message);
  }
};

export const createCommunityEvent = async (req, res) => {
  try {
    const { title, subtitle, description, location, date, timeStart, timeEnd, accessLevel, capacity } = req.body;
    const newEvent = await prisma.event.create({
      data: {
        title,
        subtitle,
        description,
        location,
        date: new Date(date),
        timeStart,
        timeEnd,
        accessLevel,
        capacity,
        category: 'GoogleDeveloperGroups', // Default untuk Community Event
        eventType: 'Community',
      },
    });
    return successResponse(res, 201, 'Community Event created', newEvent);
  } catch (error) {
    return errorResponse(res, 500, error.message);
  }
};

export const updateCommunityEvent = async (req, res) => {
  try {
    const updatedEvent = await prisma.event.update({
      where: {
        id: parseInt(req.params.id),
      },
      data: {
        ...req.body,
        date: req.body.date ? new Date(req.body.date) : undefined,
      },
    });
    if (updatedEvent.eventType !== 'Community') {
      return errorResponse(res, 404, 'Community Event not found');
    }
    return successResponse(res, 200, 'Community Event updated', updatedEvent);
  } catch (error) {
    if (error.code === 'P2025') {
      return errorResponse(res, 404, 'Community Event not found');
    }
    return errorResponse(res, 500, error.message);
  }
};

export const deleteCommunityEvent = async (req, res) => {
  try {
    const deletedEvent = await prisma.event.delete({
      where: {
        id: parseInt(req.params.id),
      },
    });
    if (deletedEvent.eventType !== 'Community') {
      return errorResponse(res, 404, 'Community Event not found');
    }
    return successResponse(res, 200, 'Community Event deleted', deletedEvent);
  } catch (error) {
    if (error.code === 'P2025') {
      return errorResponse(res, 404, 'Community Event not found');
    }
    return errorResponse(res, 500, error.message);
  }
};

export const rsvpCommunityEvent = async (req, res) => {
  try {
    const event = await prisma.event.findUnique({
      where: {
        id: parseInt(req.params.id),
      },
      include: {
        attendees: true,
      },
    });
    if (!event || event.eventType !== 'Community') {
      return errorResponse(res, 404, 'Community Event not found');
    }
    // Cek apakah user sudah RSVP
    const alreadyRSVP = event.attendees.some((att) => att.userId === req.user.id);
    if (alreadyRSVP) {
      return errorResponse(res, 400, 'You have already RSVPed for this event');
    }
    // Tambahkan RSVP user
    const updatedEvent = await prisma.event.update({
      where: {
        id: parseInt(req.params.id),
      },
      data: {
        attendees: {
          create: {
            userId: req.user.id,
            rsvp: true,
            attended: false,
          },
        },
      },
      include: {
        attendees: true,
      },
    });
    return successResponse(res, 200, 'RSVP successful', updatedEvent);
  } catch (error) {
    return errorResponse(res, 500, error.message);
  }
};

export const cancelRsvpCommunityEvent = async (req, res) => {
  try {
    const event = await prisma.event.findUnique({
      where: {
        id: parseInt(req.params.id),
      },
      include: {
        attendees: true,
      },
    });
    if (!event || event.eventType !== 'Community') {
      return errorResponse(res, 404, 'Community Event not found');
    }
    // Cek apakah user sudah RSVP
    const attendee = event.attendees.find((att) => att.userId === req.user.id);
    if (!attendee) {
      return errorResponse(res, 400, 'You have not RSVPed for this event');
    }
    // Hapus RSVP user
    await prisma.eventAttendee.delete({
      where: {
        userId_eventId: {
          userId: req.user.id,
          eventId: parseInt(req.params.id),
        },
      },
    });
    const updatedEvent = await prisma.event.findUnique({
      where: {
        id: parseInt(req.params.id),
      },
      include: {
        attendees: true,
      },
    });
    return successResponse(res, 200, 'RSVP cancelled', updatedEvent);
  } catch (error) {
    return errorResponse(res, 500, error.message);
  }
};

export const markAttendance = async (req, res) => {
  try {
    const { id, userId } = req.params;
    const event = await prisma.event.findUnique({
      where: {
        id: parseInt(id),
      },
      include: {
        attendees: true,
      },
    });
    if (!event || event.eventType !== 'Community') {
      return errorResponse(res, 404, 'Community Event not found');
    }
    // Cari attendee berdasarkan userId
    const attendee = event.attendees.find((att) => att.userId === parseInt(userId));
    if (!attendee) {
      return errorResponse(res, 404, 'Attendee not found');
    }
    // Tandai kehadiran
    const updatedEvent = await prisma.event.update({
      where: {
        id: parseInt(id),
      },
      data: {
        attendees: {
          update: {
            where: {
              userId_eventId: {
                userId: parseInt(userId),
                eventId: parseInt(id),
              },
            },
            data: {
              attended: true,
            },
          },
        },
      },
      include: {
        attendees: true,
      },
    });
    return successResponse(res, 200, 'Attendance marked', updatedEvent);
  } catch (error) {
    return errorResponse(res, 500, error.message);
  }
};